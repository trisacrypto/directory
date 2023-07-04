import { Heading } from '@chakra-ui/react';
import MemberTable from './MemberTable';
import { Trans } from '@lingui/macro';
import Loader from 'components/Loader';
import { useFetchMembers } from './hook/useFetchMembers';
const MemberPage: React.FC = () => {
  const { members, isFetchingMembers, error } = useFetchMembers();

  if (isFetchingMembers) return <Loader />;

  return (
    <>
      <Heading marginBottom="69px">
        <Trans>TRISA Member Directory</Trans>
      </Heading>
      {error && <p>error </p>}
      {members && <MemberTable data={members} />}
    </>
  );
};

export default MemberPage;
