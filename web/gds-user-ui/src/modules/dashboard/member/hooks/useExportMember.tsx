import { useSelector } from "react-redux";
import { useFetchMember } from "./useFetchMember";
import { memberSelector } from "../member.slice";
import { useState } from "react";
import { downloadMemberToCSV } from "../utils";

const useExportMember = () => {
    // Fetch the member's vasp id from the redux store
    // const { member } = useFetchMember()
    const [isLoading, setIsLoading] = useState<boolean>(false)
    const LOADING_TIMEOUT = 500;
    const exportHandler = () => {
        try{
            setIsLoading(true);
            setTimeout(() => {
                // downloadMemberToCSV(member);
                setIsLoading(false);
        }, LOADING_TIMEOUT);
        } catch (error) {
            console.error('Error exporting member', error);
        }
    };
    return {
        exportHandler,
        isLoading
    };
};

export default useExportMember;