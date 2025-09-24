import { CourseUserRoleEnum, type LoggedUserCourseDTO2, type LoggedUserDTO } from './api_types';
import { jwtDecode } from 'jwt-decode';
import { API } from '$lib/services/api.svelte';
import { localStore } from './localstore.svelte';
import { goto } from '$app/navigation';

class GlobalState {
	public reloadCounter = $state(0);

	private _accessToken = localStore<string | null>('access_token', null);
	get accessToken() {
		return this._accessToken.value;
	}
	set accessToken(val) {
		this._accessToken.value = val;

		if (val) {
			let userData = jwtDecode<{
				exp: number;
				user: LoggedUserDTO;
			}>(val);

			if (userData.exp < Date.now() / 1000) {
				API.refreshAccessToken();
			} else {
				this.loggedUser = userData.user;
			}
		} else {
			this.loggedUser = null;
			console.log('Transfering 1');
			goto('/login');
		}
	}

	private _loggedUser = localStore<LoggedUserDTO | null>('logged_user', null);
	get loggedUser() {
		return this._loggedUser.value;
	}
	set loggedUser(val) {
		this._loggedUser.value = val;

		if (val) {
			this.setCourse(val);
		} else {
			console.log('Transfering 2');
			goto('/login');
		}
	}

	private _activeCourse = $state<LoggedUserCourseDTO2 | null>(null);
	get activeCourse() {
		return this._activeCourse;
	}
	set activeCourse(val) {
		this._activeCourse = val;

		if (val) {
			if (!this.activeRole) {
				if (val.roles[0]) {
					this.activeRole = this.activeCourse?.roles[0] as CourseUserRoleEnum;
				} else {
					console.log('Transfering 3');
					goto('/login');
				}
			} else {
				if (!this.activeCourse?.roles.includes(this.activeRole)) {
					if (val.roles[0]) {
						this.activeRole = this.activeCourse?.roles[0] as CourseUserRoleEnum;
					} else {
						console.log('Transfering 4');
						goto('/login');
					}
				}
			}
		} else {
			console.log('Transfering 5');
			goto('/login');
		}
	}

	private _availableCourses = localStore<LoggedUserCourseDTO2[]>('available_courses', []);
	get availableCourses() {
		return this._availableCourses.value;
	}
	set availableCourses(val) {
		this._availableCourses.value = val;

		if (this.loggedUser) {
			this.setCourse(this.loggedUser);
		} else {
			console.log('Transfering 6');
			goto('/login');
		}
	}

	private _activeRole = localStore<CourseUserRoleEnum | undefined>('active_role', undefined);
	get activeRole() {
		return this._activeRole.value;
	}
	set activeRole(val) {
		this._activeRole.value = val;

		console.log('Transfering 7');
		if (this.activeCourse) {
			goto('/app/' + this.activeCourse.id);
		} else {
			goto('/app/');
		}
	}

	setCourse(loggedUser: LoggedUserDTO) {
		if (this.activeCourse != null) {
			const courseID = this.activeCourse.id;
			const matchedCourse = this.availableCourses.find((course) => course.id === courseID);

			if (matchedCourse) {
				this.activeCourse = matchedCourse;
			} else {
				this.activeCourse = null;
			}
		} else {
			if (loggedUser.courses.length > 0) {
				this.activeCourse = this.availableCourses[0];
			} else {
				this.activeCourse = null;
			}
		}
	}
}

export default new GlobalState();
