//
//  Server.swift
//  HAL 9000 (iOS)
//
//  Created by John Jones on 4/5/22.
//

import Foundation
import Combine
import CoreLocation
import CryptoKit

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
    private let clientId: String
    private let key: SymmetricKey
    private var subscribers = Set<AnyCancellable>()
    private let locationManager: CLLocationManager
    @Published var coordinate: Coordinate?
    @Published var messages: [MessageHolder] = []
    @Published var error: Error?
    @Published var communicating: Bool = false
    @Published var authorizationStatus: CLAuthorizationStatus
    @Published var commands: Commands = Commands(commands: [String:CommandInfo]())
    
    init(apiRoot _apiRoot: String, clientId _clientId: String, key _key: String) {
        apiRoot = _apiRoot
        clientId = _clientId
        key = SymmetricKey(data: Data(_key.utf8))
        locationManager = CLLocationManager()
        authorizationStatus = locationManager.authorizationStatus
        super.init()
        locationManager.delegate = self
        locationManager.desiredAccuracy = kCLLocationAccuracyBest
        locationManager.startUpdatingLocation()
    }
    
    override init() {
        self.apiRoot = ""
        self.clientId = ""
        self.key = SymmetricKey(size:.bits128)
        self.locationManager = CLLocationManager()
        self.authorizationStatus = locationManager.authorizationStatus
    }
    
    func ping() {
        print("ping")
        self.request(method: "GET", path: "/api/ping", body: nil as String?) { (p: Pong) in
            print("pong")
        }
    }
    
    
    func send(req: Inbound) {
        self.messages.append(MessageHolder(message: req, timestamp: Date()))
        self.request(method: "POST", path: "/api/request", body: req) { (m: Outbound) in
            self.messages.append(MessageHolder(message: m, timestamp: Date()))
            if self.messages.count > 50 {
                self.messages.removeFirst()
            }
        }
    }
    
    func getCommands() {
        self.request(method: "GET", path: "/api/commands", body: nil as String?) { (commands: Commands) in
            self.commands = commands
        }
    }
    
    func clearError() {
        self.error = nil
    }
    
    private func request<T, V>(method: String, path: String, body: T?, action: @escaping (V) -> Void) where T : Encodable, V : Decodable {
        do {
            var request = URLRequest(url: URL(string: self.apiRoot+path)!)
            request.httpMethod = method
            
            if let b = body {
                request.httpBody = try JSONEncoder().encode(b)
                request.setValue("application/json", forHTTPHeaderField: "Content-type")
            }
            
            let RFC3339DateFormatter = DateFormatter()
            RFC3339DateFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ss.SSSxxx"
            let reqTime = RFC3339DateFormatter.string(from: Date())
            
            let signature = HMAC<SHA256>.authenticationCode(for: Data("\(self.clientId):\(reqTime)".utf8), using: self.key)
            let signatureStr = Data(signature).map { String(format: "%02hhx", $0) }.joined()
            
            request.setValue(self.clientId, forHTTPHeaderField: "User-Agent")
            request.setValue(reqTime, forHTTPHeaderField: "X-Request-Time")
            request.setValue(signatureStr, forHTTPHeaderField: "X-Signature")
            
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
