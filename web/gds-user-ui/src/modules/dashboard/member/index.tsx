import { Heading, Text } from '@chakra-ui/react';
import MemberTable from './MemberTable';
import { Trans } from '@lingui/macro';
import Loader from 'components/Loader';
import { useFetchMembers } from './hook/useFetchMembers';
import Card from 'components/ui/Card';

const MemberPage: React.FC = () => {
  const { members, isFetchingMembers, error } = useFetchMembers();

  if (isFetchingMembers) return <Loader />;
  console.log('members', members);
  return (
    <>
      <Heading marginBottom="69px">
        <Trans>TRISA Member Directory</Trans>
      </Heading>
      <Card maxW="100%" marginBottom={6}>
        <Card.Body>
          <Text>
            <Trans>
              The TRISA Member Directory is for informational purposes and is meant to foster collaboration between members. It is available to verified TRISA contacts only. If you cannot view the list, please complete the registration process. You can view the member list after verification is complete.
            </Trans>
          </Text>
        </Card.Body>
      </Card>

      {error && <p>error </p>}
      {members && <MemberTable data={members} />}
    </>
  );
};

export default MemberPage;
