import { Td, Tr } from '@chakra-ui/react';

import UnverifiedMember from './UnverifiedMember';
import Loader from 'components/Loader';
import MemberTableRow from './MemberTableRow';
import { VaspType } from '../memberType';
export interface MemberTableRowsProps {
  rows: any;
  hasError: boolean;
  isLoading?: boolean;
}

const renderUnverifiedError = () => {
  return (
    <Tr>
      <Td colSpan={6}>
        <UnverifiedMember />
      </Td>
    </Tr>
  );
};

const renderLoader = () => {
  return (
    <Tr>
      <Td colSpan={6}>
        <Loader />
      </Td>
    </Tr>
  );
};

const MemberTableRows: React.FC<MemberTableRowsProps> = ({ rows, hasError, isLoading }) => {
  if (hasError) {
    return renderUnverifiedError();
  }
  if (isLoading) {
    return renderLoader();
  }
  return (
    <>
      {rows?.map((row: VaspType) => (
        <MemberTableRow row={row} key={row?.id} />
      ))}
    </>
  );
};

export default MemberTableRows;
