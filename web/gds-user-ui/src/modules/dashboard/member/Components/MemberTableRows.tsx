/* eslint-disable @typescript-eslint/no-unused-vars */
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { Td, Tr, Button, HStack, chakra } from '@chakra-ui/react';

import { BsEye } from 'react-icons/bs';

interface MemberTableRowsProps {
  rows: any;
}

const MemberTableRows: React.FC<MemberTableRowsProps> = (rows: any) => {
  console.log('rows', rows);
  return (
    <Tr>
      <>
        <Td>
          <chakra.span display="block"></chakra.span>
        </Td>
        <Td></Td>
        <Td></Td>
        <Td></Td>
        <Td></Td>
        <Td paddingY={0}></Td>
      </>
    </Tr>
  );
};

export { MemberTableRows };
