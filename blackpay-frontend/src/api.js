import axios from "axios";

export const initiatePayment = async (amount, email, currency) => {
  const response = await axios.post("http://localhost:8080/api/payment", {
    amount,
    email,
    currency,
  });

  return response.data;
};
