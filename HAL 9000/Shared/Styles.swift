//
//  Styles.swift
//  HAL 9000
//
//  Created by John Jones on 4/6/22.
//

import SwiftUI

//["Jura-Regular", "Jura-Light", "Jura-Medium", "Jura-SemiBold", "Jura-Bold"]

extension Font {
    static func halTitle() -> Font {
        return Font.custom("Jura-Bold", size: 36)
    }
    
    static func halLabel() -> Font {
        return Font.custom("Jura-SemiBold", size: 22)
    }
    
    static func halBigInput() -> Font {
        return Font.custom("Jura-Medium", size: 24)
    }
    
    static func halMedium() -> Font {
        return Font.custom("Jura-SemiBold", size: 16)
    }
    
    static func halSmall() -> Font {
        return Font.custom("Jura-Regular", size: 14)
    }
}

extension Text {
    static func hal(_ str: String) -> some View {
        return Text(str)
            .foregroundColor(Color.white)
            .textCase(Text.Case.uppercase)
    }
}

extension TextField where Label == Text {
    static func hal(_ titleKey: LocalizedStringKey, text: Binding<String>) -> some View {
        return TextField(titleKey, text: text)
            .foregroundColor(Color.white)
            .padding()
            .overlay(
                Rectangle()
                    .stroke(Color.white, lineWidth: 1)
            )
    }
    
    static func halSmall(_ titleKey: LocalizedStringKey, text: Binding<String>) -> some View {
        return TextField(titleKey, text: text)
            .foregroundColor(Color.white)
            .padding(EdgeInsets(top: 5, leading: 5, bottom: 5, trailing: 5))
            .overlay(
                Rectangle()
                    .stroke(Color.white, lineWidth: 1)
            )
    }
}

extension SecureField where Label == Text {
    static func hal(_ titleKey: LocalizedStringKey, text: Binding<String>) -> some View {
        return SecureField(titleKey, text: text)
            .foregroundColor(Color.white)
            .padding()
            .overlay(
                Rectangle()
                    .stroke(Color.white, lineWidth: 1)
            )
    }
}

extension Button where Label == Text {
    static func hal(_ titleKey: LocalizedStringKey, action: @escaping () -> Void) -> some View {
        return Button(titleKey, action: action)
            .padding(EdgeInsets(top: 15, leading: 0, bottom: 15, trailing: 0))
            .frame(maxWidth: .infinity)
            .background(Color.white)
            .foregroundColor(Color.black)
            .textCase(Text.Case.uppercase)
    }
    
    static func halSmall(_ titleKey: LocalizedStringKey, action: @escaping () -> Void) -> some View {
        return Button(titleKey, action: action)
            .padding(EdgeInsets(top: 5, leading: 5, bottom: 5, trailing: 5))
            .background(Color.white)
            .foregroundColor(Color.black)
            .textCase(Text.Case.uppercase)
    }
}
