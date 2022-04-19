//
//  ContentView.swift
//  Shared
//
//  Created by John Jones on 4/5/22.
//

import SwiftUI

struct ContentView: View {
    @ObservedObject var hal = HAL9000(apiRoot: ApiRoot, clientId: ClientId, key: Key)
    var body: some View {
        ZStack {
            Chat(hal: hal)
            hal.error != nil ? ErrorPopup(error: hal.error!) {
                hal.error = nil
            }.frame(maxHeight: .infinity, alignment: .bottom) : nil
        }.onAppear() {
            hal.ping()
        }
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
