import React, { useState, useRef, useLayoutEffect, useEffect } from 'react';
import { Box, Heading, Stack, Text, SimpleGrid, List, ListItem, Tag } from '@chakra-ui/react';
import { BUSINESS_CATEGORY, getBusinessCategiryLabel } from 'constants/basic-details';
import { getNationalIdentificationLabel } from 'constants/national-identification';
import { COUNTRIES } from 'constants/countries';
import { renderAddress } from 'utils/address-utils';
import { hasValue } from 'utils/utils';
import { t, Trans } from '@lingui/macro';

type OrganizationalDetailProps = {
  data: any;
  network?: string;
  status: string;
};

const OrganizationalDetail: React.FC<OrganizationalDetailProps> = ({ data, network, status }) => {
  const getOrgDivEl: any = document.getElementById('org') as HTMLDivElement;
  const getCntDivEl: any = document.getElementById('cnt') as HTMLDivElement;
  const [, setDivOrgHeight] = useState(getOrgDivEl);
  const [, setDivCntHeight] = useState(getCntDivEl);
  const orgRef = useRef<HTMLDivElement>(null);
  const cntRef = useRef<HTMLDivElement>(null);
  const [mainnetProfileStatus, setMainnetProfileStatus] = useState(t`pending registration`);
  const [testnetProfileStatus, setTestnetProfileStatus] = useState(t`pending registration`);
  useLayoutEffect(() => {
    setDivOrgHeight(orgRef?.current?.clientHeight || 500);
    setDivCntHeight(cntRef?.current?.clientHeight || 500);
  }, [getOrgDivEl, getCntDivEl]);

  useEffect(() => {
    if (network === 'mainnet' && status) {
      setMainnetProfileStatus(data?.organization_name);
    }
  }, [data?.organization_name, network, status]);

  useEffect(() => {
    if (network === 'testnet' && status) {
      setTestnetProfileStatus(data?.organization_name);
    }
  }, [data?.organization_name, network, status]);

  return (
    <Stack w="full">
      <Stack bg={'#E5EDF1'} justifyItems={'center'}>
        <Stack mb={5}>
          <Heading fontSize={20} pt={4} pl={8}>
            <Trans>
              {`
                Your ${network === 'mainnet' ? 'MainNet' : 'TestNet'} TRISA Organization Profile:
              `}
            </Trans>

            <Text as={'span'} pl={2} color={'#55ACD8'}>
              [{network === 'mainnet' ? mainnetProfileStatus : testnetProfileStatus}]
            </Text>
          </Heading>
        </Stack>
      </Stack>
      <SimpleGrid columns={{ base: 1, sm: 1, lg: 2 }} spacing={{ lg: 10 }}>
        <Stack border={'1px solid #eee'} p={4} mt={5} mb={8} px={8} bg={'white'} id={'org'} ref={orgRef}>
          <Box pb={2}>
            <Heading as={'h1'} fontSize={19} pb={8} pt={2}>
              <Trans>Organizational Details</Trans>
            </Heading>
            <SimpleGrid minChildWidth={"280px"} spacing="20px">
              <List>
                <ListItem fontWeight={'bold'} pb={1}>
                  <Trans>Name Identifiers</Trans>
                </ListItem>
                <ListItem>
                  {data?.entity?.name?.name_identifiers?.[0]?.legal_person_name || 'N/A'}
                </ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'} pb={1}>
                  <Trans>Organization Type</Trans>
                </ListItem>
                <ListItem>{(BUSINESS_CATEGORY as any)[data?.business_category] || 'N/A'}</ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'} pb={1}>
                  <Trans>VASP Category</Trans>
                </ListItem>
                <ListItem>
                  {data?.vasp_categories && data?.vasp_categories.length > 0
                    ? data?.vasp_categories?.map((categ: any, index: any) => {
                        return (
                          <Tag key={index} color={'white'} bg={'blue'} mr={2} mb={1} size={'lg'}>
                            {getBusinessCategiryLabel(categ)}
                          </Tag>
                        );
                      })
                    : 'N/A'}
                </ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'} pb={1}>
                  <Trans>Incorporation Date</Trans>
                </ListItem>
                <ListItem>{data?.established_on || 'N/A'}</ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'} pb={1}>
                  <Trans>Business Address</Trans>
                </ListItem>
                <ListItem>
                  {renderAddress(data?.entity?.geographic_addresses?.[0] || 'N/A')}
                </ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'} pb={1}>
                  <Trans>Identification Number</Trans>
                </ListItem>
                <ListItem>
                  {data?.entity?.national_identification?.national_identifier || 'N/A'}
                </ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'} pb={1}>
                  <Trans>Identification Type</Trans>
                </ListItem>
                <ListItem>
                  {' '}
                  {data?.entity?.national_identification?.national_identifier_type ? (
                    <Tag color={'white'} bg={'blue'} size={'lg'}>
                      {getNationalIdentificationLabel(
                        data?.entity?.national_identification?.national_identifier_type
                      )}
                    </Tag>
                  ) : (
                    'N/A'
                  )}
                </ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'} pb={1}>
                  <Trans>Country of Registration</Trans>
                </ListItem>
                <ListItem>
                  {(COUNTRIES as any)[data?.entity?.country_of_registration] || 'N/A'}
                </ListItem>
              </List>
            </SimpleGrid>
          </Box>
        </Stack>
        <Stack
          border={'1px solid #eee'}
          p={4}
          px={8}
          mt={5}
          mb={8}
          bg={'white'}
          // minHeight={divHeight}
          // boxSize={'content-box'}
          id={'cnt'}>
          <Box>
            <Heading as={'h1'} fontSize={19} pb={8} pt={2}>
              <Trans>Contacts</Trans>
            </Heading>
            <SimpleGrid spacing="20px">
              {['legal', 'technical', 'administrative', 'billing'].map((contact, index) => (
                <List key={index}>
                  <ListItem fontWeight={'bold'} pb={1}>
                    {' '}
                    {contact === 'legal'
                      ? `Compliance / ${contact.charAt(0).toUpperCase() + contact.slice(1)}`
                      : contact.charAt(0).toUpperCase() + contact.slice(1)}
                  </ListItem>
                  <ListItem>
                    {hasValue(data?.contacts?.[contact]) ? (
                      <>
                        {data?.contacts?.[contact]?.name && (
                          <Text>{data?.contacts?.[contact]?.name}</Text>
                        )}
                        {data?.contacts?.[contact]?.email && (
                          <Text>{data?.contacts?.[contact]?.email}</Text>
                        )}
                        {data?.contacts?.[contact]?.phone && (
                          <Text>{data?.contacts?.[contact]?.phone}</Text>
                        )}
                      </>
                    ) : (
                      'N/A'
                    )}
                  </ListItem>
                </List>
              ))}
            </SimpleGrid>
          </Box>
        </Stack>
      </SimpleGrid>
    </Stack>
  );
};

export default OrganizationalDetail;
