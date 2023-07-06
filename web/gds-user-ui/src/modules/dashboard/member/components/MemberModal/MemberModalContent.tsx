import { Tag, Text, Stack, SimpleGrid, List, ListItem } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import { Member } from '../../memberType';
import { BUSINESS_CATEGORY, getBusinessCategiryLabel } from 'constants/basic-details';
import { hasValue } from 'utils/utils';
interface MemberDetailProps {
  member: Member;
}

const MemberDetail = ({ member }: MemberDetailProps) => {
  return (
    <>
      <Stack w="full" pb={5}>
        <SimpleGrid minChildWidth="280px" spacing="20px" maxH={'700px'} overflowY={'scroll'}>
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
            <ListItem>
              {(BUSINESS_CATEGORY as any)[member?.summary?.business_category] || 'N/A'}
            </ListItem>
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
                <Trans>
                  {contact === 'legal'
                    ? `Compliance / ${contact.charAt(0).toUpperCase() + contact.slice(1)} Contact`
                    : contact.charAt(0).toUpperCase() + contact.slice(1) + ' Contact'}
                </Trans>
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
      </Stack>
    </>
  );
};

export default MemberDetail;
