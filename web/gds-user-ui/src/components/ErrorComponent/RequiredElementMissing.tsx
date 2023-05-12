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
      <Trans>Please make sure you have filled out all required fields.</Trans>
    </Alert>
  );
};

export default RequiredElementMissing;
