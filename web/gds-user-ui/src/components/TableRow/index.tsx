import RowItem from "./RowItem";
import { IconButton, Menu, MenuButton, MenuItem, MenuList, Tag, TagLabel, Td } from "@chakra-ui/react";
import { BsThreeDots } from "react-icons/bs";
import React, { ReactNode } from "react";

type TableRowProps<T> = {
    row: T | {[k: string]: ReactNode}
};

function TableRow<T>({ row }: TableRowProps<T>) {
    return <>
            <RowItem>
                    {
                        Object.entries(row).map(([k, v]) => (
                            <Td key={k}>{v}</Td>
                        ))
                    }
            </RowItem>
    </>;
}

export default TableRow;
