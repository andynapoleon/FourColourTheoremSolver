import "./globals.css";
import type { Metadata } from "next";
import { Anek_Gurmukhi } from "next/font/google";
import ConditionalNavBar from "./components/ConditionalNavBar";

const inter = Anek_Gurmukhi({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Map Coloring",
  description:
    "A Next.js solver application for the famous Four-Colour Map Theorem. Built collaboratively during DevelopEd 2.0 by Dev Edmonton Society.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <ConditionalNavBar />
        {children}
      </body>
    </html>
  );
}
