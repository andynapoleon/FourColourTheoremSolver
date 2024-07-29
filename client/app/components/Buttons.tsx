"use client";
import { useState } from "react";
import styles from "./styles/Buttons.module.css";
import { useRouter } from "next/navigation";

export function SignInButton() {
  const router = useRouter();

  const handleSignIn = () => {
    router.push("/login");
  };

  return (
    <button
      className={`${styles.button} ${styles.signInButton}`}
      onClick={handleSignIn}
    >
      Sign in
    </button>
  );
}

export function SignOutButton({ onSignOut }: { onSignOut: () => void }) {
  return (
    <button
      className={`${styles.button} ${styles.signOutButton}`}
      onClick={onSignOut}
    >
      Sign out
    </button>
  );
}
