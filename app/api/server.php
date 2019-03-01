<?php
/**
 * 服务器管理模块
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019-02-28
 * Time: 16:02
 */

namespace app\api;


class server
{
    private $uid,$db;
    public function __construct()
    {
        // 获取用户UID并同时判断登陆
        $this->uid = \app\Service\Verify::isLogin();
        $this->db = Data::pdo();
    }

    /**
     * 首页列表
     * @param int $pageid 分页id，lastid
     * @param null $search 搜索关键词
     * @return array
     */
    public function index($pageid = 0, $search = null)
    {
        $sql = 'select * from `server` where uid=:uid';
        $params = [
            ':uid' => $this->uid
        ];
        if (!empty($search)) {
            $sql .= ' and server_name like :search';
            $params[':search'] = "%{$search}%";
        }
        if ($pageid > 0) {
            $sql .= ' and id<:pageid';
            $params[':pageid'] = $pageid;
        }
        $sql .= ' order by id desc limit 20';
        $data = Data::pdo()->quickPrepare($sql, $params)->toArray();
        return [1, ['data' => $data]];
    }

    /**
     * 设置部署数据
     * @param int $id 应用编号，修改的时候用
     * @param string $title 应用标题
     * @param int $deploy_type 应用类型 0-无 1-git 2-zip
     * @param string $remote_url 资源地址
     * @param string $local_path 本地部署地址
     * @param string $branch git分支名
     * @param string $before_command 前置命令
     * @param string $after_command 后置命令
     * @return int
     */
    public function set($id = 0, $title = '', $deploy_type = 0, $remote_url = '', $local_path = '', $branch = '', $before_command = '', $after_command = '')
    {
        !empty($title) || Response::json(-4, null, '标题不能为空');
        in_array($deploy_type, [0,1,2]) || Response::json(-4, null, '部署类型不支持');
        $params = [
            ':title' => $title,
            ':deploy_type' => $deploy_type,
            ':remote_url' => $remote_url,
            ':local_path' => $local_path,
            ':branch' => $branch,
            ':before_command' => $before_command,
            ':after_command' => $after_command
        ];
        if ($id > 0) {
            $sql_a = [];
            foreach ($params as $k => $v) {
                $sql_a[] = '`' . substr($k, 1) . '`='. $k;
            }
            $sql_a = implode(',', $sql_a);
            $params[':id'] = $id;
            $params[':uid'] = $this->uid;
            $statement = $this->db->quickPrepare("update deploy set {$sql_a} where id=:id and uid=:uid", $params);
        } else {
            $params[':uid'] = $this->uid;
            $sql_a = $sql_b = [];
            foreach ($params as $k => $v) {
                $sql_a[] = '`' . substr($k, 1) . '`';
                $sql_b[] = $k;
            }
            $sql_a = implode(',', $sql_a);
            $sql_b = implode(',', $sql_b);
            $statement = $this->db->quickPrepare("insert into deploy({$sql_a}) values ({$sql_b});", $params);
        }
        $result = $statement->rowCount() == 1;
        $statement->closeCursor();
        return $result ? 1 : 0;
    }

    /**
     * 删除应用
     * @param int $id 应用编号
     * @return int
     */
    public function delete($id = 0)
    {
        $id > 0 || Response::json(-4, null, '应用不存在');
        $statement = Data::pdo()->quickPrepare('delete from deploy where id=:id and uid=:uid', [
            ':id' => $id,
            ':uid' => $this->uid
        ]);
        $result = $statement->rowCount() == 1;
        $statement->closeCursor();
        if ($result) {
            return 1;
        }
        return 0;
    }
}