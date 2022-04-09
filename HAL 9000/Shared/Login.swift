//
//  Logiin.swift
//  HAL 9000 (iOS)
//
//  Created by John Jones on 4/6/22.
//

import SwiftUI

struct Login: View {
    @ObservedObject var hal: HAL9000
    @State var username: String = ""
    @State var password: String = ""
    
    var body: some View {
        VStack{
            VStack{
                Text.hal("Login")
                    .padding(.bottom, 20)
                    .font(Font.halTitle())
                Text.hal("Username")
                    .font(Font.halLabel())
                    .frame(maxWidth: .infinity, alignment: .leading)
                TextField.hal("Username", text: $username)
                    .font(Font.halBigInput())
                Text.hal("Password")
                    .font(Font.halLabel())
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.top, 20)
                SecureField.hal("Password", text: $password)
                    .font(Font.halBigInput())
            }.padding(.bottom, 30)
            Button.hal("Login") {
                self.hal.login(req: LoginRequest(username: self.username, password: self.password))
            }
            .disabled(self.hal.communicating || hal.error != nil)
            .frame(maxWidth: .infinity)
            .font(Font.halLabel())
        }
        .padding()
        .frame( maxWidth: .infinity, maxHeight: .infinity)
        .background(Color.black)
    }
}

struct Logiin_Previews: PreviewProvider {
    static var previews: some View {
        Login(hal: HAL9000())
    }
}
