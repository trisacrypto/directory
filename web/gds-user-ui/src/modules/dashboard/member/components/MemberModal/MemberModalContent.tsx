import { Tag, Text, Stack, Box, SimpleGrid, List, ListItem } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import { Member } from '../../memberType';
import { getBusinessCategiryLabel } from 'constants/basic-details';
import { hasValue } from 'utils/utils';
interface MemberDetailProps {
  member: Member;
}

const MemberDetail = ({ member }: MemberDetailProps) => {
  return (
    <>
      <Stack w="full">
        <SimpleGrid minChildWidth="120px" spacing="40px">
          <Stack border={'1px solid #eee'} p={4} my={5} px={7} bg={'white'} id={'org'}>
            <Box pb={5}>
              <SimpleGrid minChildWidth="280px" spacing="40px">
                <List>
                  <ListItem fontWeight={'bold'}>
                    <Trans>Website</Trans>
                  </ListItem>
                  <ListItem>{member?.summary?.website || 'N/A'}</ListItem>
                </List>
                <List>
                  <ListItem fontWeight={'bold'}>
                    <Trans>Business Category</Trans>
                  </ListItem>
                  <ListItem>{member?.summary?.business_category || 'N/A'}</ListItem>
                </List>
                <List>
                  <ListItem fontWeight={'bold'}>
                    <Trans>Vasp Category</Trans>
                  </ListItem>
                  <ListItem>
                    {member?.summary?.vasp_categories && member?.summary?.vasp_categories.length > 0
                      ? member?.summary?.vasp_categories?.map((categ: any, index: any) => {
                          return (
                            <Tag key={index} color={'white'} bg={'blue'} mr={2} mb={1}>
                              {getBusinessCategiryLabel(categ)}
                            </Tag>
                          );
                        })
                      : 'N/A'}
                  </ListItem>
                </List>
                <List>
                  <ListItem fontWeight={'bold'}>
                    <Trans>Country of Registration</Trans>
                  </ListItem>
                  <ListItem>{member?.summary?.country || 'N/A'}</ListItem>
                </List>
                {['legal', 'technical', 'administrative'].map((contact, index) => (
                  <List key={index}>
                    <ListItem fontWeight={'bold'}>
                      {' '}
                      {contact === 'legal'
                        ? `Compliance / ${contact.charAt(0).toUpperCase() + contact.slice(1)}`
                        : contact.charAt(0).toUpperCase() + contact.slice(1)}
                    </ListItem>
                    <ListItem>
                      {hasValue(member?.contacts?.[contact]) ? (
                        <>
                          {member?.contacts?.[contact]?.name && (
                            <Text>{member?.contacts?.[contact]?.name}</Text>
                          )}
                          {member?.contacts?.[contact]?.email && (
                            <Text>{member?.contacts?.[contact]?.email}</Text>
                          )}
                          {member?.contacts?.[contact]?.phone && (
                            <Text>{member?.contacts?.[contact]?.phone}</Text>
                          )}
                        </>
                      ) : (
                        'N/A'
                      )}
                    </ListItem>
                  </List>
                ))}
                <List>
                  <ListItem fontWeight={'bold'}>
                    <Trans>TRISA Endpoint</Trans>
                  </ListItem>
                  <ListItem>{member?.summary?.endpoint || 'N/A'}</ListItem>
                </List>
                <List>
                  <ListItem fontWeight={'bold'}>
                    <Trans>Common Name</Trans>
                  </ListItem>
                  <ListItem>{member?.summary?.common_name || 'N/A'}</ListItem>
                </List>
              </SimpleGrid>
            </Box>
          </Stack>
        </SimpleGrid>
      </Stack>
    </>
  );
};

export default MemberDetail;
