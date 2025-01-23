"use client";
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import styles from "./styles/SignUpForm.module.css";
import Image from "next/image";
import Link from "next/link";

const SignUpForm: React.FC = () => {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    const apiHost = process.env.NEXT_PUBLIC_API_GATEWAY_URL;
    if (!apiHost) {
      throw new Error("API host is not defined in the environment variables");
    }

    try {
      const response = await fetch(`${apiHost}/api/v1/auth/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify({ name, email, password }),
      });

      if (response.ok) {
        const data = await response.json();
        console.log("Sign-up successful:", data);
        router.push("/login");
      } else {
        const data = await response.json();
        setError(data.error);
      }
    } catch (error) {
      console.error("Error during sign-up:", error);
      setError("An error occurred during sign-up");
    }
  };

  return (
    <div className={styles.signupContainer}>
      <form onSubmit={handleSubmit} className={styles.signupForm}>
        <Image
          src="/logo.png"
          alt="Cartoon Logo"
          width={100}
          height={100}
          className={styles.logo}
        />
        <h2 className={styles.formTitle}>Create an Account</h2>
        <div className={styles.inputGroup}>
          <label htmlFor="name" className={styles.label}>
            Name
          </label>
          <input
            type="text"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className={styles.inputField}
            required
          />
        </div>
        <div className={styles.inputGroup}>
          <label htmlFor="email" className={styles.label}>
            Email
          </label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className={styles.inputField}
            required
          />
        </div>
        <div className={styles.inputGroup}>
          <label htmlFor="password" className={styles.label}>
            Password
          </label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className={styles.inputField}
            required
          />
        </div>
        {error && <p className={styles.error}>{error}</p>}
        <button type="submit" className={styles.submitButton}>
          Sign Up
        </button>
        <p className={styles.loginLink}>
          Already have an account? <Link href="/login">Log in</Link>
        </p>
      </form>
    </div>
  );
};

export default SignUpForm;
