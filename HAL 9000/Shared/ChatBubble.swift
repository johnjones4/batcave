//
//  ChatBubble.swift
//  HAL 9000 (iOS)
//
//  Created by John Jones on 4/5/22.
//

import SwiftUI

struct ChatBubble<Content>: View where Content: View {
    let message: MessageHolder
    let content: () -> Content
    let color: Color
    init(message: MessageHolder, @ViewBuilder content: @escaping () -> Content) {
        self.content = content
        self.message = message
        self.color = message.message is Inbound ? Color.red : Color.blue
    }
    
    var body: some View {
        HStack(spacing: 0 ) {
            content()
                .padding(.all, 15)
                .foregroundColor(Color.white)
                .background(color)
                .clipShape(Rectangle())
        }
    }
}
