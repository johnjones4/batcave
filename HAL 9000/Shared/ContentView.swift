//
//  ContentView.swift
//  Shared
//
//  Created by John Jones on 4/5/22.
//

import SwiftUI

struct ContentView: View {
    @ObservedObject var hal = try! HAL9000(apiRoot: "https://hal9000.frontdoor.johnjonesfour.com")
    var body: some View {
        ZStack {
            (hal.token != nil && !(hal.token?.isExpired ?? true)) ? AnyView(Chat(hal: hal)) : AnyView(Login(hal: hal))
            hal.error != nil ? ErrorPopup(error: hal.error!) {
                hal.error = nil
            }.frame(maxHeight: .infinity, alignment: .bottom) : nil
        }
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
