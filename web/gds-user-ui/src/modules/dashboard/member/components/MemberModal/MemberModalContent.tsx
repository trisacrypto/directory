import { Tag, Text, Stack, SimpleGrid, List, ListItem } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import { Member } from '../../memberType';
import { getBusinessCategiryLabel } from 'constants/basic-details';
import { hasValue } from 'utils/utils';
import { getBusinessCategory } from '../../utils';
interface MemberDetailProps {
  member: Member;
}

const MemberDetail = ({ member }: MemberDetailProps) => {
  return (
    <>
      <Stack w="full" pb={5}>
        <SimpleGrid spacing="20px" maxH={'700px'} overflowY={'auto'}>
          <List>
            <ListItem fontWeight={'bold'}>
              <Trans>Website</Trans>
            </ListItem>
            <ListItem>{member?.data?.summary?.website || 'N/A'}</ListItem>
          </List>
          <List>
            <ListItem fontWeight={'bold'}>
              <Trans>Business Category</Trans>
            </ListItem>
            <ListItem>
              {getBusinessCategory(member?.data?.summary?.business_category as any) || 'N/A'}
            </ListItem>
          </List>
          <List>
            <ListItem fontWeight={'bold'}>
              <Trans>VASP Category</Trans>
            </ListItem>
            <ListItem>
              {member?.data?.summary?.vasp_categories && member?.data?.summary?.vasp_categories.length > 0
                ? member?.data?.summary?.vasp_categories?.map((categ: any, index: any) => {
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
            <ListItem>{member?.data?.summary?.country || 'N/A'}</ListItem>
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
                {hasValue(member?.data?.contacts?.[contact]) ? (
                  <>
                    {member?.data?.contacts?.[contact]?.name && (
                      <Text>{member?.data?.contacts?.[contact]?.name}</Text>
                    )}
                    {member?.data?.contacts?.[contact]?.email && (
                      <Text>{member?.data?.contacts?.[contact]?.email}</Text>
                    )}
                    {member?.data?.contacts?.[contact]?.phone && (
                      <Text>{member?.data?.contacts?.[contact]?.phone}</Text>
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
            <ListItem>{member?.data?.summary?.endpoint || 'N/A'}</ListItem>
          </List>
          <List>
            <ListItem fontWeight={'bold'}>
              <Trans>Common Name</Trans>
            </ListItem>
            <ListItem>{member?.data?.summary?.common_name || 'N/A'}</ListItem>
          </List>
        </SimpleGrid>
      </Stack>
    </>
  );
};

export default MemberDetail;
