import React, { useState } from "react";
import image from "../assets/images/login/login.png";
import axios from "axios"; 
import { useLocation } from "react-router-dom";

export default function OtpForm() {
  const [otp, setOtp] = useState(["", "", "", "", "", ""]);
 const location = useLocation();
 const email = location.state?.email;
  const handleChange = (e, index) => {
    const value = e.target.value;
    if (!/^[0-9]?$/.test(value)) return; // Only allow single digit numbers

    const newOtp = [...otp];
    newOtp[index] = value;
    setOtp(newOtp);

    // Auto-focus next input
    if (value && index < otp.length - 1) {
      document.getElementById(`otp-${index + 1}`).focus();
    }
  };

  const handleKeyDown = (e, index) => {
    if (e.key === "Backspace" && !otp[index] && index > 0) {
      document.getElementById(`otp-${index - 1}`).focus();
    }
  };

  const handleSubmit = async () => {
    const fullOtp = otp.join("");
    console.log("OTP Submitted:", fullOtp);

     try {
      //  Send POST request using Axios
      const response = await axios.post("http://localhost:9000/verify-otp", {
        email: email,
        otp: fullOtp,
      });

      console.log("OTP verifed");
      alert("Login successful!");

      // Navigate to OTP or next page
      navigate("/");
    } catch (error) {
      console.error("Login error:", error);
      alert(
        error.response?.data?.message ||
          "Login failed. Please check your credentials."
      );
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="flex flex-row items-center justify-between w-full max-w-7xl">
        {/* Left Image Section */}
        <div className="flex justify-center items-center w-1/2 -ml-20">
          <img
            src={image}
            alt="OTP Illustration"
            className="w-[650px] h-[650px] object-cover"
          />
        </div>

        {/* Right OTP Form Section */}
        <div className="w-[500px] bg-gradient-to-br from-white via-blue-50 to-blue-200 rounded-2xl shadow-2xl p-12 flex flex-col justify-center">
          <h1 className="text-4xl font-bold text-blue-700 mb-6 text-center">
            Verify OTP
          </h1>
          <p className="text-gray-600 text-center mb-8">
            Enter the 6-digit code sent to your email or phone
          </p>

          {/* OTP Inputs */}
          <div className="flex justify-center gap-4 mb-8">
            {otp.map((digit, index) => (
              <input
                key={index}
                id={`otp-${index}`}
                type="text"
                maxLength="1"
                value={digit}
                onChange={(e) => handleChange(e, index)}
                onKeyDown={(e) => handleKeyDown(e, index)}
                className="w-12 h-12 text-center text-2xl font-semibold text-blue-700 bg-white border border-blue-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
              />
            ))}
          </div>

          {/* Buttons */}
          <div className="grid grid-cols-2 gap-4">
            <button
              onClick={handleSubmit}
              className="px-6 py-3 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition shadow-md"
            >
              Verify
            </button>
            <button
              onClick={() => console.log("Resend OTP clicked")}
              className="px-6 py-3 bg-white text-blue-600 font-semibold rounded-lg border-2 border-blue-600 hover:bg-blue-50 transition"
            >
              Resend OTP
            </button>
          </div>

          <p className="text-center text-sm text-gray-600 mt-6">
            Didn’t receive the code?{" "}
            <span
              onClick={() => console.log("Try another method")}
              className="text-blue-600 font-medium cursor-pointer hover:underline"
            >
              Try another method
            </span>
          </p>
        </div>
      </div>
    </div>
  );
}
