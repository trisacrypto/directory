import { FiCompass } from 'react-icons/fi';
import { FaRegLightbulb, FaBook } from 'react-icons/fa';
import { HiOutlineUserGroup } from 'react-icons/hi';
import { BiCertification } from 'react-icons/bi';

const Menu = [
  {
    title: 'Overview',
    icon: FiCompass,
    activated: true,
    path: '/dashboard/overview'
  },
  {
    title: 'Certificate Management',
    icon: BiCertification,
    activated: true,
    path: '/dashboard/certificate-management'
  },
  {
    title: 'Technical Resources',
    icon: FaRegLightbulb
  },
  {
    title: 'Collaborators',
    icon: HiOutlineUserGroup,
    path: '/dashboard/overview'
  },
  {
    title: 'Member Directory',
    icon: FaBook
  }
];

export default Menu;
