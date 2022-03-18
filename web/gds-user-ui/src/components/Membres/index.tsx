import React from 'react';
import {
  Tr,
  Thead,
  Table,
  Flex,
  Tbody,
  Th,
  Stack,
  Heading,
  useColorModeValue,
  FormControl,
  Input,
  Button,
  InputRightElement,
  MenuGroup
} from '@chakra-ui/react';
import MembersRow from 'components/Tables/MembersRow';
import Card from 'components/ui/Card';
import { SearchIcon } from '@chakra-ui/icons';
const Members = ({ datas }: any) => {
  const textColor = useColorModeValue('#858585', 'white');
  return (
    <Card>
      <Stack p={4} mb={5}>
        <Heading fontSize="20px" fontWeight="bold" pb=".5rem">
          Members
        </Heading>
      </Stack>
      <Stack px={'60px'}>
        <Flex>
          <FormControl color={'gray.500'} width={'35%'}>
            <Input size="md" pr="4.5rem" type={'gray.100'} placeholder="Search by Member name" />
            <InputRightElement width="2.5rem" color={'black'}>
              <Button h="2.5rem" size="sm" onClick={(e) => {}}>
                <SearchIcon />
              </Button>
            </InputRightElement>
          </FormControl>
        </Flex>
        <Table
          color={textColor}
          width={'100%'}
          sx={{
            borderCollapse: 'separate',
            borderSpacing: '0 10px',
            Th: {
              textTransform: 'capitalize',
              color: '#858585',
              fontWeight: 'bold',
              borderBottom: 'none',
              fontSize: '0.9rem'
            }
          }}>
          <Thead>
            <Tr>
              <Th>Member Name</Th>
              <Th>TestNet</Th>
              <Th>MainNet</Th>
              <Th pr={0}>Member Details</Th>
            </Tr>
          </Thead>
          <Tbody>
            {datas.map((row: any) => {
              return (
                <MembersRow
                  key={row.name}
                  name={row.name}
                  isTestNet={row.isTestNet}
                  isMainNet={row.isMainNet}
                  memberId={row.memberId}
                />
              );
            })}
          </Tbody>
        </Table>
      </Stack>
    </Card>
  );
};
Members.defaultProps = {
  datas: [
    {
      name: 'John Doe',
      isTestNet: true, // true or false
      isMainNet: true,
      memberId: '123456789'
    }
  ]
};

export default Members;
