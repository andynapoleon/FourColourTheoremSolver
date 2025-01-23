"use client";

import Link from "next/link";
import styles from "./styles/NavMenu.module.css";
import Image from "next/image";
import { ProfileButton, SignInButton, SignOutButton } from "./Buttons";
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

  const handleSignOut = async () => {
    const apiHost = process.env.NEXT_PUBLIC_API_GATEWAY_URL;
    if (!apiHost) {
      throw new Error("API host is not defined in the environment variables");
    }

    console.log("Token:", localStorage.getItem("token"));

    try {
      const response = await fetch(`${apiHost}/api/v1/auth/logout`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });

      if (response.ok) {
        console.log("Sign-out successful", response);
        localStorage.removeItem("token");
        localStorage.removeItem("name");
        setIsAuthenticated(false);
        router.push("/login"); // Redirect to login page after sign out
      } else {
        console.error("Error during sign-out:", response.statusText);
      }
    } catch (error) {
      console.error("Error during sign-out:", error);
    }
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
        {isAuthenticated ? (
          <>
            <ProfileButton />
            <SignOutButton onSignOut={handleSignOut} />
          </>
        ) : (
          <SignInButton />
        )}
      </div>
    </nav>
  );
}
