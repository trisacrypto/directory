import { Td, Tr, Button, HStack, chakra, Tag } from '@chakra-ui/react';
import { BsEye } from 'react-icons/bs';
import { formatIsoDate } from 'utils/formate-date';
import { getVapsNetwork } from '../utils';
import { VaspType } from '../memberType';

const MemberTableRow: React.FC<{ row: VaspType }> = ({ row }) => {
  return (
    <Tr key={row?.id}>
      <Td>
        <chakra.span display="block">{row?.name}</chakra.span>
      </Td>
      <Td>{formatIsoDate(row?.first_listed)}</Td>
      <Td>{formatIsoDate(row?.last_updated)}</Td>
      <Td>{getVapsNetwork(row?.registered_directory)}</Td>
      <Td>
        <Tag bg="green.400" color="white">
          {row?.status}
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
  );
};

export default MemberTableRow;
