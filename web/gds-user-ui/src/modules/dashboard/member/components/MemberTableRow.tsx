import { Td, Tr, chakra, Tag } from '@chakra-ui/react';

import { formatIsoDate } from 'utils/formate-date';
import { getVapsNetwork, getVaspStatus } from '../utils';
import { VaspType } from '../memberType';
import ShowMemberModal from '../components/MemberModal';

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
          {getVaspStatus(row?.status)}
        </Tag>
      </Td>
      <Td paddingY={0}>
        <ShowMemberModal memberId={row?.id} />
      </Td>
    </Tr>
  );
};

export default MemberTableRow;
