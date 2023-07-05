import { Td, Tr, Button, HStack, chakra, Tag } from '@chakra-ui/react';
import { BsEye } from 'react-icons/bs';
import { formatIsoDate } from 'utils/formate-date';
import { VaspDirectoryEnum } from '../memberType';

interface MemberTableRowsProps {
  rows: any;
}

const MemberTableRows: React.FC<MemberTableRowsProps> = (rows: any) => {
  // for now we are just displaying the unverified member component
  // with the status check story we will check if the member is verified or not and display the appropriate component
  return rows.rows.data.map((row: any) => (
    <Tr key={row.id}>
      <Td>
        <chakra.span display="block">{row.name}</chakra.span>
      </Td>
      <Td>{formatIsoDate(row.first_listed)}</Td>
      <Td>{formatIsoDate(row.last_updated)}</Td>
      {row.registered_directory === VaspDirectoryEnum.MAINNET ||
      row.registered_directory === VaspDirectoryEnum.MAINNETDEV ? (
        <Td>MainNet</Td>
      ) : (
        <Td>TestNet</Td>
      )}
      <Td>
        <Tag bg="green.400" color="white">
          {row.status}
        </Tag>
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
    </Tr>
  ));
};

export { MemberTableRows };
