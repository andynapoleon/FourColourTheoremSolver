"use client";

import { usePathname } from "next/navigation";
import NavBar from "./NavBar";

export default function ConditionalNavBar() {
  const pathname = usePathname();

  // Don't render NavBar on login page
  if (pathname === "/login") {
    return null;
  }

  return <NavBar />;
}
