//
//  Coordinate.swift
//  HAL 9000 (iOS)
//
//  Created by John Jones on 4/5/22.
//

import Foundation

protocol Message: Encodable, Decodable {
    var body: String { get }
}

struct MessageHolder {
    let message: Message
    let timestamp: Date
}

struct Coordinate: Encodable, Decodable {
    let latitude: Double
    let longitude: Double
}

struct Inbound: Message {
    let body: String
    let location: Coordinate
}

struct Outbound: Message {
    let body: String
    let media: String
    let url: String
}

struct LoginRequest: Encodable {
    let username: String
    let password: String
}

struct Token: Encodable, Decodable {
    let token: String
    let user: String
    let expiration: String
    
    var expirationDate: Date {
        let RFC3339DateFormatter = DateFormatter()
        RFC3339DateFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ss.SSSxxx"
        return RFC3339DateFormatter.date(from: expiration)!
    }
    
    var isExpired: Bool {
        return expirationDate <= Date()
    }
}

struct Commands: Decodable {
    let commands: [String: String]
    
    func suggest(partial: String) -> [String] {
        if partial.count == 0 || partial[partial.startIndex] != "/" {
            return [String]()
        }
        let partialLc = partial.lowercased()
        return commands
            .map({ k, _ in
                return "/" + k
            })
            .filter { command in
                return command.lowercased().starts(with: partialLc)
            }
    }
}

struct ErrorResponse: Decodable {
    let error: String
    let code: Int?
    
    var description : String {
        return error
    }
}
