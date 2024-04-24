import Link from "next/link";
import styles from "./NavMenu.module.css";
import Image from "next/image";
import { SignInButton, SignOutButton } from "./components/buttons";

export default function NavBar() {
  return (
    <nav className={styles.nav}>
      <Link href={"/"}>
        <Image
          src="/logo.png" // Route of the image file
          width={50}
          height={30}
          alt="Map Coloring Logo"
        />
      </Link>
      <h1>The Best Map Coloring App in the World!</h1>
    </nav>
  );
}
