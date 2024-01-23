import {
  Alert,
  AlertIcon,
  ListItem,
  UnorderedList,
  AlertDescription,
  Box,
  Text,
  HStack
} from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
interface RequiredElementMissingProps {
  elementKey?: number;
  errorFields?: any[];
}

const RequiredElementMissing = ({ errorFields }: RequiredElementMissingProps) => {
  return (
    <Alert
      status="error"
      flexDirection="column"
      textAlign="left"
      alignItems="left"
      justifyContent="left"
      variant="subtle"
      borderRadius="lg"
      my={4}>
      <HStack>
        <AlertIcon />
        <Trans>There is an error or missing data in this section. Please edit the section.</Trans>
      </HStack>

      <Box ml="35" mt="2">
        <AlertDescription>
          <UnorderedList>
            {errorFields?.map((errorField: any, index) => {
              return (
                <ListItem key={index}>
                  <Text as="span" fontWeight="bold">
                    {errorField?.field}
                    {': '}
                  </Text>
                  <Text as="span">{errorField?.error}</Text>
                </ListItem>
              );
            })}
          </UnorderedList>
        </AlertDescription>
      </Box>
    </Alert>
  );
};

export default RequiredElementMissing;
