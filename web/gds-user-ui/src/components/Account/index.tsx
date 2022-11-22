import { Avatar, Box, HStack, Text } from '@chakra-ui/react';
import { Organization } from 'modules/dashboard/organization/organizationType';
import { BsChevronRight } from 'react-icons/bs';
import { Link } from 'react-router-dom';

type AccountProps = {
  src?: string;
  alt?: string;
} & Partial<Organization>;

export const Account = ({ domain, name, src }: AccountProps) => (
  <Link to="/#" role="group">
    <HStack
      direction="row"
      w="100%"
      _groupHover={{ bg: '#BEE3F81d' }}
      padding={3}
      transition="background 400ms">
      <HStack width="100%" alignItems="center">
        <HStack>
          <Avatar bg="blue" name={name} src={src} fontWeight={700} />
          <Box>
            <Text fontWeight={700} data-testid="vaspName">
              {name}
            </Text>
            <Text data-testid="vaspDomain">{domain}</Text>
          </Box>
        </HStack>
      </HStack>
      <BsChevronRight />
    </HStack>
  </Link>
);
