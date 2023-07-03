import { FiCompass } from 'react-icons/fi';
import { FaRegLightbulb, FaRegMoneyBillAlt } from 'react-icons/fa';
import { BiCertification } from 'react-icons/bi';
import { IconType } from 'react-icons';
import { CheckCircleIcon } from '@chakra-ui/icons';
import { ComponentWithAs, IconProps } from '@chakra-ui/react';
import { HiUserGroup } from 'react-icons/hi';
import { GoFileDirectory } from 'react-icons/go';
import { t } from '@lingui/macro';
// import CertificateManagementIcon from 'assets/certificate-management.svg';

type Menu = {
  title: string;
  icon?: IconType | ComponentWithAs<'svg', IconProps>;
  activated?: boolean;
  path?: string;
  children?: Menu[];
  isExternalLink?: boolean;
};

const MENU: Menu[] = [
  {
    title: t`Overview`,
    icon: FiCompass,
    activated: true,
    path: '/dashboard/overview'
  },
  {
    title: t`Certificate Management`,
    icon: BiCertification,
    activated: true,
    children: [
      {
        title: t`Certificate Registration`,
        icon: CheckCircleIcon,
        path: '/dashboard/certificate/registration',
        activated: true
      },
      {
        title: t`Certificate Inventory`,
        icon: FaRegMoneyBillAlt,
        path: '/dashboard/certificate/inventory',
        activated: !!(process.env.REACT_APP_ENABLE_CERT_MANAGEMENT_FEAT === 'true')
      }
    ]
  },
  {
    title: t`Collaborators`,
    activated: true,
    icon: HiUserGroup,
    path: '/dashboard/collaborators'
  },
  {
    title: t`Member Directory`,
    activated: true,
    icon: GoFileDirectory,
    path: '/dashboard/member'
  },
  {
    title: t`Technical Resources`,
    icon: FaRegLightbulb
  }
];

export default MENU;
