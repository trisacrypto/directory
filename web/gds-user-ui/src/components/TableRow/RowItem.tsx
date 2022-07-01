import React from "react";
import { Tr } from "@chakra-ui/react";

const RowItem: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    return (
        <Tr
            border="1px solid #23A7E0"
            borderRadius={100}
            css={{
                'td:first-child': {
                    border: '1px solid #23A7E0',
                    borderRight: 'none',
                    borderTopLeftRadius: 100,
                    borderBottomLeftRadius: 100
                },
                'td:last-child': {
                    border: '1px solid #23A7E0',
                    borderLeft: 'none',
                    borderTopRightRadius: 100,
                    borderBottomRightRadius: 100,
                    textAlign: 'center'
                },
                'td:not(:first-child):not(:last-child)': {
                    borderTop: '1px solid #23A7E0',
                    borderBottom: '1px solid #23A7E0'
                }
            }}>
            {children}
        </Tr>
    );
};

export default RowItem;
