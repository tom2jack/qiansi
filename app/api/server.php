<?php
/**
 * 服务器管理模块
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019-02-28
 * Time: 16:02
 */

namespace app\api;

use Respect\Validation\Validator as v;
use zhimiao\Data;
use zhimiao\Response;
use zhimiao\Utils;

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
     * 添加/修改服务器配置
     * @param int $id 服务器注册编号
     * @param string $server_name 服务器备注名
     * @param string $address 服务器地址
     * @param string $api_secret 接口密钥
     * @param string $api_rule 接口权限规则
     * @return array|int
     */
    public function set($id = 0, $server_name = '', $address = '', $api_secret = '', $api_rule = '')
    {
        !empty($server_name) || Response::json(-4, null, '标题不能为空');
        !empty($api_secret) || Response::json(-4, null, '接口密钥必填');
        if (!v::ip()->validate($address) && !v::domain()->validate($address)) {
            return [-4, null, '地址格式错误'];
        }
        $params = [
            ':server_name' => $server_name,
            ':address' => $address,
            ':api_secret' => $api_secret,
            ':api_rule' => $api_rule,
        ];
        if ($id > 0) {
            $sql_a = [];
            foreach ($params as $k => $v) {
                $sql_a[] = '`' . substr($k, 1) . '`='. $k;
            }
            $sql_a = implode(',', $sql_a);
            $params[':id'] = $id;
            $params[':uid'] = $this->uid;
            $statement = $this->db->quickPrepare("update `server` set {$sql_a} where id=:id and uid=:uid", $params);
        } else {
            $params[':uid'] = $this->uid;
            $sql_a = $sql_b = [];
            foreach ($params as $k => $v) {
                $sql_a[] = '`' . substr($k, 1) . '`';
                $sql_b[] = $k;
            }
            $sql_a = implode(',', $sql_a);
            $sql_b = implode(',', $sql_b);
            $statement = $this->db->quickPrepare("insert into `server`({$sql_a}) values ({$sql_b});", $params);
        }
        $result = $statement->rowCount() == 1;
        $statement->closeCursor();
        return $result ? 1 : 0;
    }

    /**
     * 删除服务器
     * @param int $id 应用编号
     * @return int
     */
    public function delete($id = 0)
    {
        $id > 0 || Response::json(-4, null, '应用不存在');
        $statement = Data::pdo()->quickPrepare('delete from `server` where id=:id and uid=:uid', [
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

    /**
     * 校验服务器状态
     * @param int $id
     */
    public function check($id = 0)
    {
        $id > 0 || Response::json(-4, null, '服务器不存在');
        $info = $this->db->quickPrepare('select * from server where id=:id and uid=:uid', [
            ':id' => $id,
            ':uid' => $this->uid
        ])->getOnce();
        if ($info === false) {
            return [-5, null, '查无此服务器'];
        }
        if (empty($info['api_secret'])) {
            return [-5, null, '请先配置服务器API密钥'];
        }
        $client = new \GuzzleHttp\Client(['timeout' => 1]);
        $check_data = \app\Service\Utils::encrypt('zhimiao.api.check', $info['api_secret']);
        $url = 'http://'. $info['address']. ':1314/check';
        $url = 'http://localhost:1305/index.php?a=miao';
        try {
            $result = $client->post($url, [
                'body' => $check_data
            ]);
            if ($result->getStatusCode() != 200) {
                return [-6, null, '请求失败'];
            }
            $data = \app\Service\Utils::decrypt($result->getBody(), $info['api_secret']);
            if ($data != 'succ') {
                return [-7, null, '验证失败'];
            }
        } catch (\Exception $e) {
            return [-6, null, '请求失败'];
        } catch (\Error $e2) {
            return [-6, null, '请求失败'];
        }
        return 1;
    }
}