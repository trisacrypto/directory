import { Td, Tr, Button, HStack, chakra, Tag } from '@chakra-ui/react';

import UnverifiedMember from '../UnverifiedMember';
import { BsEye } from 'react-icons/bs';
import { VaspDirectoryEnum } from '../memberType';

interface MemberTableRowsProps {
  rows: any;
}

const MemberTableRows: React.FC<MemberTableRowsProps> = (rows: any) => {
  const isMainnet = rows.registered_directory === VaspDirectoryEnum.MAINNET;
  // for now we are just displaying the unverified member component
  // with the status check story we will check if the member is verified or not and display the appropriate component
  return (
    <Tr>
      {rows.length > 0 ? (
        rows.rows.data.map((row: any) => (
          <>
          <Td key={row.id}>
            <chakra.span display="block">{row.name}</chakra.span>
          </Td>
          <Td>{row.first_listed}</Td>
          <Td>{row.last_updated}</Td>
          <Td>{isMainnet ? <span>MainNet</span> : <span>TestNet</span>}</Td>
          <Td>
            <Tag bg="green.400" color="white">{row.status}</Tag>
            </Td>
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
        ))
      ) : (
        <Td colSpan={6}>
          <UnverifiedMember />
        </Td>
      )}
    </Tr>
  );
};

export { MemberTableRows };
