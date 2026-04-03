graph TD
    A[Start: Create Lead] --> B[Enter Name, Number, Check-in/out]
    B --> C[Select Bag Count]
    
    C --> D{Are Bags >= 1?}
    D -- No --> E[Disable Promo Code Field]
    D -- Yes --> F[Enable Promo Code Field]
    
    F --> G[User Clicks Apply Promo]
    G --> H[Call Validate Coupon API]
    
    H --> I{Is Service Selected?}
    I -- No --> J[Auto-add Default Service Count = 1]
    I -- Yes --> K[Calculate Discount Only]
    
    J --> L[Call Rentcal API]
    K --> L[Call Rentcal API]
    
    L --> M[Proceed to Payment]
    M --> N[Payment Successful]
    N --> O[Show Summary Page]
    O --> P[Generate Check-In QR Code]
    
    P --> Q[Check-in Complete & Visible on Live Dashboard]
    
    Q --> R[User Initiates Checkout]
    R --> S{Exceeded Checkout Time?}
    S -- Yes --> T[Apply Exceed Charges: ₹15 / hr / bag]
    T --> U[Settle Exceed Amount]
    U --> V[Checkout Complete]
    S -- No --> V[Checkout Complete]
