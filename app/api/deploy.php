<?php
/**
 * 部署服务模块
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019-02-18
 * Time: 19:31
 */

namespace app\api;

use Respect\Validation\Validator as v;
use zhimiao\Data;
use zhimiao\Request;

class deploy
{
    private $uid;
    public function __construct()
    {
        // 获取用户UID并同时判断登陆
        $this->uid = \app\Service\Verify::isLogin();
    }

    /**
     * 首页列表
     * @param int $pageid 分页id，lastid
     * @param null $search 搜索关键词
     * @return array
     */
    public function index($pageid = 0, $search = null)
    {
        $sql = 'select * from deploy where uid=:uid';
        $params = [
            ':uid' => $this->uid
        ];
        if (!empty($search)) {
            $sql .= ' and title like :search';
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
     */
    public function set()
    {
        $id = Request::post('id', 0);
        $title = Request::post('title');
        $deploy_type = Request::post('deploy_type');
        $remote_url = Request::post('remote_url', '');
        $local_path = Request::post('local_path', '');
        $branch = Request::post('branch', '');
        $before_command = Request::post('before_command', '');
        $after_command = Request::post('after_command', '');
        if (!v::notEmpty()->validate($title)) {
            return [-4, null, '标题不能为空'];
        }

    }
}