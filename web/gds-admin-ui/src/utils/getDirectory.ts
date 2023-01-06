import { DIRECTORY_NAME } from "@/constants";
import { isTestNet } from ".";

export default function getDirectory(){
    return isTestNet() ? DIRECTORY_NAME.TRISATEST : DIRECTORY_NAME.VASP_DIRECTORY
}