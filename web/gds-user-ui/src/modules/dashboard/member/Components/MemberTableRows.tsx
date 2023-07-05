import { Td, Tr, Button, HStack, chakra } from '@chakra-ui/react';

import UnverifiedMember from '../UnverifiedMember';
import { BsEye } from 'react-icons/bs';

interface MemberTableRowsProps {
  rows: any;
}

const MemberTableRows: React.FC<MemberTableRowsProps> = (rows: any) => {
  // for now we are just displaying the unverified member component
  // with the status check story we will check if the member is verified or not and display the appropriate component
  return (
    <Tr>
      {rows.length > 0 ? (
        <>
          <Td>
            <chakra.span display="block"></chakra.span>
            <chakra.span display="block" fontSize="sm" color="gray.700"></chakra.span>
          </Td>
          <Td></Td>
          <Td></Td>
          <Td></Td>
          <Td></Td>
          <Td paddingY={0}>
            <HStack width="100%" justifyContent="center" alignItems="center">
              <Button
                color="blue"
                as={'a'}
                href={``}
                bg={'transparent'}
                _hover={{
                  bg: 'transparent'
                }}
                _focus={{
                  bg: 'transparent'
                }}>
                <BsEye fontSize="24px" />
              </Button>
            </HStack>
          </Td>
        </>
      ) : (
        <Td colSpan={6}>
          <UnverifiedMember />
        </Td>
      )}
    </Tr>
  );
};

export { MemberTableRows };
