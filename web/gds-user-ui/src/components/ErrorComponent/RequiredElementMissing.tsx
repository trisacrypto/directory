import { Alert, AlertIcon } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
interface RequiredElementMissingProps {
  elementKey?: number;
}

const RequiredElementMissing = ({ elementKey }: RequiredElementMissingProps) => {
  console.log('k', elementKey);

  return (
    <Alert status="error" borderRadius="lg" my={4}>
      <AlertIcon />
      <Trans>There is an error or missing data in this section. Please edit the section.</Trans>
    </Alert>
  );
};

export default RequiredElementMissing;
