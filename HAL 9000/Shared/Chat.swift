//
//  Chat.swift
//  HAL 9000 (iOS)
//
//  Created by John Jones on 4/5/22.
//

import SwiftUI

struct Chat: View {
    @ObservedObject var hal: HAL9000
    @State var text: String = ""
    @State var recommendations: [String] = [String]()
    
    var body: some View {
        switch hal.authorizationStatus {
        case .notDetermined:
            VStack{}
                .background(Color.black)
                .onAppear() {
                    hal.requestPermission()
                }
        case .restricted:
            Text("Location use is restricted.")
        case .denied:
            Text("The app does not have location permissions. Please enable them in settings.")
        case .authorizedAlways, .authorizedWhenInUse:
            GeometryReader { geo in
                VStack {
                    CustomScrollView(scrollToEnd: true) {
                        LazyVStack {
                            ForEach(0..<hal.messages.count, id:\.self) { index in
                                AnyView(ChatBubble(message: hal.messages[index]) {
                                    VStack{
                                        Text(hal.messages[index].message.body)
                                            .font(Font.halMedium())
                                            .frame(maxWidth: .infinity, alignment: .leading)
                                        Text(hal.messages[index].timestamp.formatted())
                                            .font(Font.halSmall())
                                            .frame(maxWidth: .infinity, alignment: .leading)
                                    }
                                })
                                if let halMessage = hal.messages[index].message as? Outbound {
                                    if halMessage.media != "" {
                                        AsyncImage(url: URL(string: halMessage.media)) { image in
                                            image
                                              .resizable()
                                              .aspectRatio(contentMode: .fit)
                                        } placeholder: {
                                            Color.gray
                                        }
                                    }
                                }
                            }
                        }
                    }.padding(.top)
                    self.recommendations.isEmpty ? nil : Button.hal(LocalizedStringKey(self.recommendations.first!), action: {
                        send()
                    })
                    HStack {
                        ZStack {
                            TextField.halSmall("Message", text: $text)
                            .onSubmit {
                                send()
                            }
                            .onChange(of: text, perform: { newValue in
                                self.recommendations = self.hal.commands.suggest(partial: newValue)
                            })
                                .font(Font.halSmall())
                        }
                        
                        Button.halSmall("Send") {
                            send()
                        }
                        .font(Font.halSmall())
                        .frame(maxHeight: .infinity)
                    }
                    .padding(EdgeInsets(top: 10, leading: 10, bottom: 10, trailing: 10))
                    .frame(maxHeight: 50)
                }
                .background(Color.black)
            }
            .onAppear() {
                self.hal.getCommands()
            }
        default:
            Text("Unexpected status")
        }
    }
    
    func send() {
        var sendText = text
        if !self.recommendations.isEmpty {
            let firstSpace = sendText.firstIndex(of: " ")
            let body = firstSpace == nil ? "" : " " + sendText[firstSpace!...]
            sendText = self.recommendations.first! + body
        }
        hal.send(req: Inbound(
            body: sendText,
            location: hal.coordinate ?? Coordinate(latitude: 0, longitude: 0)
        ))
        text = ""
    }
}
