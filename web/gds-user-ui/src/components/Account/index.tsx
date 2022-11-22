import { Avatar, Box, HStack, Text } from '@chakra-ui/react';
import { BsChevronRight } from 'react-icons/bs';
import { Link } from 'react-router-dom';

type AccountProps = {
  vaspName: string;
  username: string;
  src?: string;
  alt?: string;
};

export const Account = ({ vaspName, username, src }: AccountProps) => (
  <Link to="/#" role="group">
    <HStack
      direction="row"
      w="100%"
      _groupHover={{ bg: '#BEE3F81d' }}
      padding={3}
      transition="background 400ms">
      <HStack width="100%" alignItems="center">
        <HStack>
          <Avatar bg="blue" name={username} src={src} fontWeight={700} />
          <Box>
            <Text fontWeight={700} data-testid="vaspName">
              {vaspName}
            </Text>
            <Text data-testid="username">{username}</Text>
          </Box>
        </HStack>
      </HStack>
      <BsChevronRight />
    </HStack>
  </Link>
);
