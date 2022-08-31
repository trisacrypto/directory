import React from 'react';
import { Stack, Box, Text, Heading, UnorderedList, ListItem, HStack } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

enum DetailsCardEnum {
  ORG = 'org',
  CERT = 'cert'
}
interface DetailsCardProps {
  type: string;
  data: any;
}

const DetailsCard = ({ type, data }: DetailsCardProps) => {
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      height={248}
      maxWidth={473}
      fontSize={18}
      p={5}
      mt={10}
      px={5}>
      <Stack>
        {type === DetailsCardEnum.ORG ? (
          <Stack>
            <Heading fontSize={20}>
              <Trans id="Organizational Details">Organizational Details</Trans>
            </Heading>
            <UnorderedList p={5} mt={10} px={5}>
              <ListItem>
                <HStack justifyContent={'space-between'}>
                  <Text>
                    <Trans id="TRISA Member ID:">TRISA Member ID:</Trans>
                  </Text>
                  <Text></Text>
                </HStack>
              </ListItem>
              <ListItem>
                <HStack justifyContent={'space-between'}>
                  <Text>
                    <Trans id="TRISA Verification:">TRISA Verification:</Trans>
                  </Text>
                  <Text></Text>
                </HStack>
              </ListItem>
              <ListItem>
                <HStack justifyContent={'space-between'}>
                  <Text>
                    <Trans id="Country:">Country:</Trans>
                  </Text>
                  <Text></Text>
                </HStack>
              </ListItem>
            </UnorderedList>
          </Stack>
        ) : (
          <Stack>
            <Heading fontSize={20}>
              <Trans id="Certificate Details">Certificate Details</Trans>
            </Heading>
            <UnorderedList p={5} mt={10} px={5}>
              <ListItem>
                <HStack justifyContent={'space-between'}>
                  <Text>
                    <Trans id="Organization:">Organization:</Trans>
                  </Text>
                  <Text></Text>
                </HStack>
              </ListItem>
              <ListItem>
                <HStack justifyContent={'space-between'}>
                  <Text>
                    <Trans id="Issue Date:">Issue Date:</Trans>
                  </Text>
                  <Text></Text>
                </HStack>
              </ListItem>
              <ListItem>
                <HStack justifyContent={'space-between'}>
                  <Text>
                    <Trans id="Expiry Date:">Expiry Date:</Trans>
                  </Text>
                  <Text></Text>
                </HStack>
              </ListItem>
              <ListItem>
                <HStack justifyContent={'space-between'}>
                  <Text>
                    <Trans id="TRISA Identity Signature:">TRISA Identity Signature:</Trans>
                  </Text>
                  <Text></Text>
                </HStack>
              </ListItem>
            </UnorderedList>
          </Stack>
        )}
      </Stack>
    </Box>
  );
};
DetailsCard.defaultProps = {
  type: 'org',
  data: []
};

export default DetailsCard;
