import { Heading } from '@chakra-ui/react';
import MemberTable from './MemberTable';
import { Trans } from '@lingui/macro';
import Loader from 'components/Loader';
import { useFetchMembers } from './hooks/useFetchMembers';
import DirectoryNotification from './components/DirectoryNotification';

const MemberPage: React.FC = () => {
  const { members, isFetchingMembers, error } = useFetchMembers();

  if (isFetchingMembers) return <Loader />;
  console.log('members', members);
  return (
    <>
      <Heading marginBottom="69px">
        <Trans>TRISA Member Directory</Trans>
      </Heading>
      <DirectoryNotification />

      {error && <p>error </p>}
      {members && <MemberTable data={members} />}
    </>
  );
};

export default MemberPage;
