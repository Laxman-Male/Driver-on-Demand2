graph TD
    A[Start: Create Lead] --> B[Enter Name, Number, Check-in/out]
    B --> C[Select Bag Count]
    
C --&gt; D{Are Bags &gt;= 1?}
D -- No --&gt; E[Disable Promo Code Field]
D -- Yes --&gt; F[Enable Promo Code Field]

F --&gt; G[User Clicks Apply Promo]
G --&gt; H[Call Validate Coupon API]

H --&gt; I{Is Service Selected?}
I -- No --&gt; J[Auto-add Default Service Count = 1]
I -- Yes --&gt; K[Calculate Discount Only]

J --&gt; L[Call Rentcal API]
K --&gt; L[Call Rentcal API]

L --&gt; M[Proceed to Payment]
M --&gt; N[Payment Successful]
N --&gt; O[Show Summary Page]
O --&gt; P[Generate Check-In QR Code]

P --&gt; Q[Check-in Complete &amp; Visible on Live Dashboard]

Q --&gt; R[User Initiates Checkout]
R --&gt; S{Exceeded Checkout Time?}
S -- Yes --&gt; T[Apply Exceed Charges: ₹15 / hr / bag]
T --&gt; U[Settle Exceed Amount]
U --&gt; V[Checkout Complete]
S -- No --&gt; V[Checkout Complete]
