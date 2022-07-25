import { v4 as uuidv4 } from 'uuid';
import { UserInfo } from './utils/hex64';

const default_password = 'cd180755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368';
const default_new_password_1 = '00080755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368';
const default_new_password_2 = '11180755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368';
const default_new_password_3 = '22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368';
const random_user_email = uuidv4() + '.test@test.test';

export interface Shared {
	default_password: string;
	default_password_corrupted: string;
	default_new_password_1: string;
	default_new_password_corrupted_1: string;
	default_new_password_2: string;
	default_new_password_corrupted_2: string;
	default_new_password_3: string;
	default_new_password_corrupted_3: string;
	default_password_token_corrupted: string;
	random_user_email: string;
	random_user_email_corrupted: string;
	random_user_token: string | undefined;
	random_user_token_corrupted: string;
	random_user_userInfo: UserInfo | undefined;
}

export const shared = <Shared>{
	default_password: default_password,
	default_password_corrupted: default_password.slice(0, -1),
	default_new_password_1: default_new_password_1,
	default_new_password_corrupted_1: default_new_password_1.slice(0, -1),
	default_new_password_2: default_new_password_2,
	default_new_password_corrupted_2: default_new_password_2.slice(0, -1),
	default_new_password_3: default_new_password_3,
	default_new_password_corrupted_3: default_new_password_3.slice(0, -1),
	default_password_token_corrupted: 'default_password_token_corrupted',
	random_user_email: random_user_email,
	random_user_email_corrupted: random_user_email + '@',
	random_user_token: undefined,
	random_user_token_corrupted: 'random_user_token_corrupted',
	random_user_userInfo: undefined
};
