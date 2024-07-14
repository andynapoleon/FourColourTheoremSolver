import "./globals.css";
import type { Metadata } from "next";
import { Anek_Gurmukhi } from "next/font/google";
import ConditionalNavBar from "./components/ConditionalNavBar";

const inter = Anek_Gurmukhi({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Four-Colour Theorem Map Coloring - Great App!",
  description:
    "A solver application for the famous Four-Colour Map Theorem. Initially built collaboratively during DevelopEd 2.0 Hackathon (2023) and later further developed.",
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
