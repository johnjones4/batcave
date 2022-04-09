//
//  ErrorPopup.swift
//  HAL 9000
//
//  Created by John Jones on 4/6/22.
//

import SwiftUI

struct ErrorPopup: View {
    let error: Error
    let dismiss: (() -> Void)
    
    var body: some View {
        Group {
            Text(error.localizedDescription)
                .padding()
        }
        .frame(maxWidth: .infinity, alignment: .bottom)
        .background(Color.red)
        .foregroundColor(Color.white)
        .padding()
        .gesture(DragGesture(minimumDistance: 0, coordinateSpace: .local)
                            .onEnded({ value in
                                if value.translation.height > 0 {
                                    dismiss()
                                }
                            }))
    }
}

struct ErrorPopup_Previews: PreviewProvider {
    static var previews: some View {
        ErrorPopup(error: NSError(domain: "test", code: 0, userInfo: [
            NSLocalizedDescriptionKey: "Test error"
        ])) {
            print("dismiss")
        }
    }
}
