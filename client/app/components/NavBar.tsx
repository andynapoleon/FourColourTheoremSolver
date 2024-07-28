"use client";

import Link from "next/link";
import styles from "./styles/NavMenu.module.css";
import Image from "next/image";
import { SignInButton, SignOutButton } from "./Buttons";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

export default function NavBar() {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const router = useRouter();

  useEffect(() => {
    const checkAuthStatus = () => {
      const token = localStorage.getItem("token");
      setIsAuthenticated(!!token);
    };

    checkAuthStatus();
    // Add an event listener to check auth status when local storage changes
    window.addEventListener("storage", checkAuthStatus);

    return () => {
      window.removeEventListener("storage", checkAuthStatus);
    };
  }, []);

  const handleSignOut = () => {
    localStorage.removeItem("token");
    setIsAuthenticated(false);
    router.push("/login"); // Redirect to home page after sign out
  };

  return (
    <nav className={styles.nav}>
      <div className={styles.navItem}>
        <Link href={"/"}>
          <Image
            src="/logo.png"
            width={50}
            height={30}
            alt="Map Coloring Logo"
          />
        </Link>
      </div>
      <div className={styles.navItem}>
        <h1>The Best Map Coloring App in the World!</h1>
      </div>
      <div className={`${styles.navItem} ${styles.authButtons}`}>
        {isAuthenticated && (
          <Link href="/profile" className={styles.profileButton}>
            Profile
          </Link>
        )}
        {isAuthenticated ? (
          <SignOutButton onSignOut={handleSignOut} />
        ) : (
          <SignInButton />
        )}
      </div>
    </nav>
  );
}
