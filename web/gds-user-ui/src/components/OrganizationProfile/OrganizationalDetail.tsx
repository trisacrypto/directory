import React, { useState, useEffect, useRef, useLayoutEffect } from 'react';
import * as Sentry from '@sentry/react';
import {
  Box,
  Heading,
  VStack,
  Flex,
  Input,
  Stack,
  Text,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  SimpleGrid,
  List,
  ListItem,
  Table,
  Tbody,
  Tr,
  Td,
  HStack,
  Tag
} from '@chakra-ui/react';
import { BUSINESS_CATEGORY, getBusinessCategiryLabel } from 'constants/basic-details';
import { getNameIdentiferTypeLabel } from 'constants/name-identifiers';
import { getNationalIdentificationLabel } from 'constants/national-identification';
import { COUNTRIES } from 'constants/countries';
import { addressType } from 'constants/address';
import { renderAddress } from 'utils/address-utils';
import { hasValue } from 'utils/utils';
type OrganizationalDetailProps = {
  data: any;
};

const OrganizationalDetail: React.FC<OrganizationalDetailProps> = ({ data }) => {
  const getOrgDivEl: any = document.getElementById('org') as HTMLDivElement;
  const getCntDivEl: any = document.getElementById('cnt') as HTMLDivElement;
  const [divOrgHeight, setDivOrgHeight] = useState(getOrgDivEl);
  const [divCntHeight, setDivCntHeight] = useState(getCntDivEl);
  const orgRef = useRef<HTMLDivElement>(null);
  const cntRef = useRef<HTMLDivElement>(null);
  useLayoutEffect(() => {
    setDivOrgHeight(orgRef?.current?.clientHeight || 500);
    setDivCntHeight(cntRef?.current?.clientHeight || 500);
  }, [getOrgDivEl, getCntDivEl]);
  // keep the same div container to have a good experience
  const divHeight = divOrgHeight >= divCntHeight ? divOrgHeight : divCntHeight;
  return (
    <Stack py={5} w="full">
      <Stack bg={'#E5EDF1'} h="55px" justifyItems={'center'} p={4}>
        <Stack mb={5}>
          <Heading fontSize={20}>
            TRISA Organization Profile:{' '}
            <Text as={'span'} color={'blue'}>
              [pending registration]
            </Text>
          </Heading>
        </Stack>
      </Stack>
      <SimpleGrid minChildWidth="120px" spacing="40px">
        <Stack
          border={'1px solid #eee'}
          p={4}
          my={5}
          px={7}
          bg={'white'}
          // minHeight={divHeight}
          id={'org'}
          // boxSize={'border-box'}
          ref={orgRef}>
          <Box pb={5}>
            <Heading as={'h1'} fontSize={19} pb={10} pt={4}>
              {' '}
              Organizational Details{' '}
            </Heading>
            <SimpleGrid minChildWidth="280px" spacing="40px">
              <List>
                <ListItem fontWeight={'bold'}>Name Identifiers</ListItem>
                <ListItem>
                  {data?.entity?.name?.name_identifiers?.[0].legal_person_name || 'N/A'}
                </ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'}>Organization Type</ListItem>
                <ListItem>{(BUSINESS_CATEGORY as any)[data.business_category] || 'N/A'}</ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'}>VASP Category</ListItem>
                <ListItem>
                  {' '}
                  {data?.vasp_categories && data?.vasp_categories.length
                    ? data?.vasp_categories?.map((categ: any) => {
                        return (
                          <Tag key={categ} color={'white'} bg={'blue'} mr={2} mb={1}>
                            {getBusinessCategiryLabel(categ)}
                          </Tag>
                        );
                      })
                    : 'N/A'}
                </ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'}>Incorporation Date</ListItem>
                <ListItem>{data?.established_on || 'N/A'}</ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'}>Business Address</ListItem>
                <ListItem>
                  {' '}
                  {renderAddress(data?.entity?.geographic_addresses?.[0] || 'N/A')}
                </ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'}>Identification Number </ListItem>
                <ListItem>
                  {data?.entity?.national_identification?.national_identifier || 'N/A'}
                </ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'}>Identification Type</ListItem>
                <ListItem>
                  {' '}
                  <Tag color={'white'} bg={'blue'} size={'lg'}>
                    {getNationalIdentificationLabel(
                      data?.entity?.national_identification?.national_identifier_type
                    )}
                  </Tag>
                </ListItem>
              </List>
              <List>
                <ListItem fontWeight={'bold'}>Country of Registration </ListItem>
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
          px={5}
          my={5}
          bg={'white'}
          // minHeight={divHeight}
          // boxSize={'content-box'}
          id={'cnt'}>
          <Box>
            <Heading as={'h1'} fontSize={19} pb={10} pt={4}>
              {' '}
              Contacts{' '}
            </Heading>
            <SimpleGrid minChildWidth="360px" spacing="40px">
              {['legal', 'technical', 'administrative', 'billing'].map((contact, index) => (
                <>
                  <List>
                    <ListItem fontWeight={'bold'}>
                      {' '}
                      {contact === 'legal'
                        ? `Compliance / ${contact.charAt(0).toUpperCase() + contact.slice(1)}`
                        : contact.charAt(0).toUpperCase() + contact.slice(1)}
                    </ListItem>
                    <ListItem>
                      {hasValue(data.contacts?.[contact]) ? (
                        <>
                          {data.contacts?.[contact]?.name && (
                            <Text>{data.contacts?.[contact]?.name}</Text>
                          )}
                          {data.contacts?.[contact]?.email && (
                            <Text>{data.contacts?.[contact]?.email}</Text>
                          )}
                          {data.contacts?.[contact]?.phone && (
                            <Text>{data.contacts?.[contact]?.phone}</Text>
                          )}
                        </>
                      ) : (
                        'N/A'
                      )}
                    </ListItem>
                  </List>
                </>
              ))}
            </SimpleGrid>
          </Box>
        </Stack>
      </SimpleGrid>
    </Stack>
  );
};

export default OrganizationalDetail;
