import { http } from '@/utils/http'

export interface TeamMember {
  uid: string
  levelName: string
  registeredAt: string
}

export interface TeamPushRes {
  list: TeamMember[]
  page: number
  pageSize: number
  pageCount: number
  totalCount: number
}

/** 当前用户直推成员 */
export function getTeamPush(params?: { page?: number; pageSize?: number }) {
  return http.get<TeamPushRes>('/users/team/push', params)
}
