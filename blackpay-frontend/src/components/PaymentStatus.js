import React from "react";

function PaymentStatus({ paymentLink, status }) {
  return (
    <div>
      <p>{status}</p>
      {paymentLink && (
        <div>
          <a href={paymentLink} target="_blank" rel="noopener noreferrer">
            Complete Payment
          </a>
        </div>
      )}
    </div>
  );
}

export default PaymentStatus;
