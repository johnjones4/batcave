//
//  Server.swift
//  HAL 9000 (iOS)
//
//  Created by John Jones on 4/5/22.
//

import Foundation
import Combine
import CoreLocation

enum ServerError: Error {
    case noToken
    case tokenExpired
    case badResponse
    case response(response: ErrorResponse)
}

extension ServerError: LocalizedError {
    var errorDescription: String? {
        switch self {
        case .noToken:
            return "No login token present"
        case .tokenExpired:
            return "Login token is expired"
        case .badResponse:
            return "Bad server response"
        case .response(let response):
            return "Error \(response.code ?? 0): \(response.error)"
        }
    }
}

class HAL9000: NSObject, ObservableObject {
    private let apiRoot: String
    private var subscribers = Set<AnyCancellable>()
    private let locationManager: CLLocationManager
    @Published var coordinate: Coordinate?
    @Published var token: Token?
    @Published var messages: [MessageHolder] = []
    @Published var error: Error?
    @Published var communicating: Bool = false
    @Published var authorizationStatus: CLAuthorizationStatus
    @Published var commands: Commands = Commands(commands: [String:String]())
    
    init(apiRoot _apiRoot: String) throws {
        apiRoot = _apiRoot
        locationManager = CLLocationManager()
        authorizationStatus = locationManager.authorizationStatus
        super.init()
        try loadToken()
        locationManager.delegate = self
        locationManager.desiredAccuracy = kCLLocationAccuracyBest
        locationManager.startUpdatingLocation()
    }
    
    override init() {
        self.apiRoot = ""
        self.locationManager = CLLocationManager()
        self.authorizationStatus = locationManager.authorizationStatus
    }
    
    func login(req: LoginRequest) {
        self.request(method: "POST", path: "/api/login", authenticate: false, body: req) { (token: Token) in
            self.token = token
            do {
                try self.saveToken()
            } catch {
                self.error = error
            }
        }
    }
    
    func send(req: Inbound) {
        self.messages.append(MessageHolder(message: req, timestamp: Date()))
        self.request(method: "POST", path: "/api/request", authenticate: true, body: req) { (m: Outbound) in
            self.messages.append(MessageHolder(message: m, timestamp: Date()))
            if self.messages.count > 50 {
                self.messages.removeFirst()
            }
        }
    }
    
    func getCommands() {
        self.request(method: "GET", path: "/api/commands", authenticate: true, body: nil as String?) { (commands: Commands) in
            self.commands = commands
        }
    }
    
    func clearError() {
        self.error = nil
    }
    
    private var tokenStorageURL: URL {
        let documentDirectory = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask)[0]
        return documentDirectory.appendingPathComponent("token.json")
    }
    
    private func loadToken() throws {
        let url = self.tokenStorageURL
        if !FileManager.default.fileExists(atPath: url.path) {
            return
        }
        guard let data = try? Data(contentsOf: url) else { return }
        let decoder = JSONDecoder()
        self.token = try decoder.decode(Token.self, from: data)
    }
    
    private func saveToken() throws {
        let encoder = JSONEncoder()
        guard let data = try? encoder.encode(self.token) else { return }
        try data.write(to: self.tokenStorageURL)
    }
    
    private func request<T, V>(method: String, path: String, authenticate: Bool, body: T?, action: @escaping (V) -> Void) where T : Encodable, V : Decodable {
        do {
            var request = URLRequest(url: URL(string: self.apiRoot+path)!)
            if let b = body {
                request.httpBody = try JSONEncoder().encode(b)
                request.setValue("application/json", forHTTPHeaderField: "Content-type")
            }
            if authenticate {
                guard let token = self.token else {
                    throw ServerError.noToken
                }
                if token.isExpired {
                    throw ServerError.tokenExpired
                }
                request.setValue(token.token, forHTTPHeaderField: "Authorization")
            }
            request.httpMethod = method
            self.communicating = true
            URLSession.shared.dataTaskPublisher(for: request)
                .tryMap{ (data, response) -> Data in
                    guard let urlResponse = response as? HTTPURLResponse else {
                        throw ServerError.badResponse
                    }
                    if urlResponse.statusCode >= 300 {
                        print(String(data: data, encoding: .utf8))
                        let err = try JSONDecoder().decode(ErrorResponse.self, from: data)
                        throw ServerError.response(response: err)
                    }
                    return data
                }
                .decode(type: V.self, decoder: JSONDecoder())
                .eraseToAnyPublisher()
                .receive(on: DispatchQueue.main)
                .sink(
                    receiveCompletion: { (completion) in
                        self.communicating = false
                        switch completion {
                        case .finished:
                            break
                        case .failure(let error):
                            self.error = error
                        }
                    },
                    receiveValue: {
                        self.communicating = false
                        action($0)
                    }
                )
                .store(in: &subscribers)
        } catch {
            self.error = error
        }
    }
}

extension HAL9000: CLLocationManagerDelegate {
    func requestPermission() {
        locationManager.requestWhenInUseAuthorization()
    }

    func locationManagerDidChangeAuthorization(_ manager: CLLocationManager) {
         authorizationStatus = manager.authorizationStatus
    }
    
    func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {
        coordinate = Coordinate(
            latitude: locations.first?.coordinate.latitude ?? 0.0,
            longitude: locations.first?.coordinate.longitude ?? 0.0
        )
    }
}
